export const getJSON = r => {

  if (r.status === 500) {
    throw new Error("internal server error");
  } else if (r.status === 401) {
    throw new Error("not authorized");
  } else if (r.status === 404) {
    throw new Error("not found");
  } 

  if (r.status === 204) {
    return;
  }

  return r.json()
};
