import React, { useCallback, useEffect, useState } from "react";
import { Router, Switch, Route } from "react-router-dom";
import AppPage from "./components/AppPage";
import CreatePage from "./components/CreatePage";
import Login from "./components/Login";
import Navbar from "./components/common/Navbar";
import history from "./globalHistory";

const AppRouter = () => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();

  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [tokenUpgraded, setTokenUpgraded] = useState(false);

  const loadAccounts = async () => {
    setLoadingAccounts(true);
    const response = await fetch("/api/v1/accounts");
    if (response.status !== 200) {
      // Somehow did not get any account back.
    }
    const accounts = await response.json();
    setAccounts(accounts.resources);
    setLoadingAccounts(false);
  };

  useEffect(() => {
    // You can't have an async function as an effect argument.
    loadAccounts();
  }, []);

  const handleAccountChosen = useCallback(async ({ selectedItem }) => {
    setSelectedAccountID(selectedItem.metadata.guid);
    setTokenUpgraded(false);
    setHasChosenAccount(true);
    const { status } = await fetch("/api/v1/authenticate/account", {
      method: "POST",
      body: JSON.stringify({
        id: selectedItem.metadata.guid,
      }),
    });
    if (status === 200) {
      setTokenUpgraded(true);
    }
  }, []);
  return (
    <Router history={history}>
      <Switch>
        <Route path="/login" exact component={Login} />
        <Route path="/create" exact>
        <Navbar
            isLoaded={!isLoadingAccounts}
            items={accounts}
            accountSelected={handleAccountChosen}
          />
          <CreatePage />
        </Route>
        <Route path="/schedule" exact>
          <Navbar
              isLoaded={!isLoadingAccounts}
              items={accounts}
              accountSelected={handleAccountChosen}
            />
        </Route>
        <Route path="/" exact>
          <Navbar
            isLoaded={!isLoadingAccounts}
            items={accounts}
            accountSelected={handleAccountChosen}
          />
          <AppPage
            hasChosenAccount={hasChosenAccount}
            tokenUpgraded={tokenUpgraded}
            accountID={accountID}
          />
        </Route>
      </Switch>
    </Router>
  );
};

export default AppRouter;
