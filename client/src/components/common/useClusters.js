import { useReducer, useEffect } from "react";
import produce from "immer";

const WAIT_FOR_ALL = false;

// Takes an array of objects and tranforms it into a map of objects, with ID
// being the key and the object being the value.
// e.g.
// [{ id: 'a1', x: 'hello' }, { id: 'b2', x: 'world' }] =>
// {
//   a1: { id: 'a1', x: 'hello' },
//   b2: { id: 'b2', x: 'world' }
// }
const arrayToMap = (arr) =>
  arr.reduce((acc, cur) => ({ ...acc, [cur.id]: cur }), {});

const grab = async (url, options) => {
  const response = await fetch(url, options);
  if (response.status !== 200) {
    throw Error();
  }
  const data = await response.json();
  return data;
};

function clusterReducer(state, action) {
  switch (action.type) {
    case "FETCH_INIT":
      return {
        ...state,
        isLoading: true,
        isError: false,
      };

    case "FETCH_SUCCESS":
      return {
        ...state,
        isLoading: false,
        isError: false,
        data: action.payload,
      };

    case "FETCH_ERROR":
      return {
        ...state,
        isLoading: false,
        isError: true,
      };

    case "UPDATE_TAG": {
      const nextState = produce(state.data, (draftState) => {
        draftState[action.id].tags = action.tags.map((t) => t.name);
      });
      return {
        ...state,
        data: nextState,
      };
    }

    case "UPDATE_ALL_TAGS": {
      const nextState = produce(state.data, (draftState) => {
        action.tags.forEach((t) => {
          if (t) {
            draftState[t.id].tags = t.tags.map((t) => t.name);
          }
        });
      });
      return {
        ...state,
        data: nextState,
      };
    }

    case "UPDATE_COST":
      const nextState = produce(state.data, (draftState) => {
        draftState[action.id].cost = action.cost.bill;
      });
      return {
        ...state,
        data: nextState,
      };

    case "UPDATE_ALL_COSTS": {
      const nextState = produce(state.data, (draftState) => {
        action.cost.forEach((c) => {
          if (c) {
            draftState[c.id].cost = c.cost.bill;
          }
        });
      });
      return {
        ...state,
        data: nextState,
      };
    }

    default:
      throw new Error();
  }
}

const useClusters = (accountID) => {
  const [state, dispatch] = useReducer(clusterReducer, {
    isLoading: false,
    isError: false,
    data: [],
  });

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;
    let cancelled = false;
    const loadData = async () => {
      dispatch({ type: "FETCH_INIT" });
      try {
        const _clusters = await grab("/api/v1/clusters", { signal });

        if (!cancelled) {
          const clusters = arrayToMap(_clusters);
          dispatch({ type: "FETCH_SUCCESS", payload: clusters });

          const tagsPromises = Object.keys(clusters).map(async (id) => {
            try {
              const _tags = await grab("/api/v1/clusters/gettag", {
                signal,
                method: "POST",
                body: JSON.stringify({
                  crn: clusters[id].crn,
                }),
              });

              const tags = _tags.items;

              if (!WAIT_FOR_ALL && !cancelled) {
                dispatch({
                  type: "UPDATE_TAG",
                  id: id,
                  tags: tags,
                });
              }
              return { id: id, tags: tags };
            } catch {
              return undefined;
            }
          });

          if (WAIT_FOR_ALL) {
            Promise.all(tagsPromises).then((tags) => {
              if (!cancelled) {
                dispatch({
                  type: "UPDATE_ALL_TAGS",
                  tags: tags,
                });
              }
            });
          }

          const costPromises = Object.keys(clusters).map(async (id) => {
            try {
              const cost = await grab("/api/v1/billing", {
                signal,
                method: "POST",
                body: JSON.stringify({
                  crn: clusters[id].crn,
                  accountID: accountID,
                  clusterID: id,
                }),
              });

              if (!WAIT_FOR_ALL && !cancelled) {
                dispatch({
                  type: "UPDATE_COST",
                  id: id,
                  cost: cost,
                });
              }
              return { id: id, cost: cost };
            } catch {
              return undefined;
            }
          });
          if (WAIT_FOR_ALL) {
            Promise.all(costPromises).then((cost) => {
              if (!cancelled) {
                dispatch({
                  type: "UPDATE_ALL_COSTS",
                  cost: cost,
                });
              }
            });
          }
        }
      } catch {
        if (!cancelled) {
          dispatch({ type: "FETCH_ERROR" });
        }
      }
    };

    loadData();

    return () => {
      cancelled = true;
      controller.abort();
    };
  }, [accountID]);

  return [
    state,
    {
      deleteClusters: () => {},
      deleteTag: () => {},
      setTag: () => {},
    },
  ];
};

export default useClusters;
