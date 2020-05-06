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
      }
      const accounts = await response.json();
      setAccounts(accounts.resources);
      setLoadingAccounts(false);

      // const selectedAccount =
      //   localStorage.getItem("accountID") ||
      //   accounts.resources[0].metadata.guid;

      // setAccountStuff(selectedAccount);
    };
    loadAccounts();
  }, [handleAccountChosen, setAccountStuff]);

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
          <CreatePage hasChosenAccount={hasChosenAccount} accountID={accountID}/>
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
