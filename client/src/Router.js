import React, { useCallback, useEffect, useState, Suspense } from "react";
import { Router, Switch, Route } from "react-router-dom";
import { Loading } from "carbon-components-react";
import Navbar from "./common/Navbar";
import history from "./globalHistory";

const AppPage = React.lazy(() => import("./pages/clusters/AppPage"));
const CreatePage = React.lazy(() => import("./pages/create/CreatePage"));
const SchedulePage = React.lazy(() => import("./pages/schedule/SchedulePage"));
const SettingsPage = React.lazy(() => import("./pages/settings/SettingsPage"));
const Login = React.lazy(() => import("./Login"));

const HolderThing = (props) => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();

  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [tokenUpgraded, setTokenUpgraded] = useState(false);

  const setAccountStuff = useCallback(async (guid) => {
    localStorage.setItem("accountID", guid);
    setSelectedAccountID(guid);
    setTokenUpgraded(false);
    setHasChosenAccount(true);
    const { status } = await fetch("/api/v1/authenticate/account", {
      method: "POST",
      body: JSON.stringify({
        id: guid,
      }),
    });
    if (status === 200) {
      setTokenUpgraded(true);
    }
  }, []);

  const handleAccountChosen = useCallback(
    async ({ selectedItem }) => {
      setAccountStuff(selectedItem.metadata.guid);
    },
    [setAccountStuff]
  );

  useEffect(() => {
    const loadAccounts = async () => {
      setLoadingAccounts(true);
      const response = await fetch("/api/v1/accounts");
      if (response.status !== 200) {
        // Somehow did not get any account back.
        return;
      }
      const accounts = await response.json();
      setAccounts(accounts.resources);
      setLoadingAccounts(false);
    };
    loadAccounts();
  }, []);

  return (
    <>
      <Navbar
        path={props.location.pathname}
        isLoaded={!isLoadingAccounts}
        items={accounts}
        accountSelected={handleAccountChosen}
      />
      <Route path="/create" exact>
        <CreatePage
          tokenUpgraded={tokenUpgraded}
          hasChosenAccount={hasChosenAccount}
          accountID={accountID}
        />
      </Route>
      <Route path="/schedule" exact>
        <SchedulePage
          tokenUpgraded={tokenUpgraded}
          hasChosenAccount={hasChosenAccount}
          accountID={accountID}
        />
      </Route>
      <Route path="/settings" exact>
        <SettingsPage
          tokenUpgraded={tokenUpgraded}
          hasChosenAccount={hasChosenAccount}
          accountID={accountID}
        />
      </Route>
      <Route path="/" exact>
        <AppPage
          hasChosenAccount={hasChosenAccount}
          tokenUpgraded={tokenUpgraded}
          accountID={accountID}
        />
      </Route>
    </>
  );
};

const AppRouter = () => {
  useEffect(() => {
    fetch("/api/v1/login").then(({ status }) => {
      if (status !== 200) {
        history.push("/login");
      }
    });
  }, []);
  //style={path === “create” ? styles.activeItem : styles.item}

  return (
    <Router history={history}>
      <Suspense fallback={<Loading />}>
        <Switch>
          <Route path="/login" exact component={Login} />
          <Route path="/" component={HolderThing} />
        </Switch>
      </Suspense>
    </Router>
  );
};

export default AppRouter;
