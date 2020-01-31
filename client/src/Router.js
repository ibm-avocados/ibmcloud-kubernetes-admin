import React, {useCallback, useEffect, useState} from "react";
import { Router, Switch, Route } from "react-router-dom";
import AppPage from "./components/AppPage";
import Login from './components/Login';
import { createBrowserHistory } from "history";
import Navbar from "./components/common/Navbar";
import { getJSON } from "./fetchUtil";
import history from './globalHistory'

const DumPage = () => {
  return <div>Hello</div>;
};

const AppRouter = () => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();


  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [tokenUpgraded, setTokenUpgraded] = useState(false);

  const loadAccounts = async () => {
    setLoadingAccounts(true);
    const accounts = await fetch("/api/v1/accounts").then(getJSON);
    setAccounts(accounts.resources);
    setLoadingAccounts(false);
  };

  useEffect(() => {
    // You can't have an async function as an effect argument.
    loadAccounts();
  }, []);

  const handleAccountChosen = useCallback(
    async ({ selectedItem }) => {
      setSelectedAccountID(selectedItem.metadata.guid);
      setTokenUpgraded(false);
      setHasChosenAccount(true);
      const { status } = await fetch("/api/v1/authenticate/account", {
        method: "POST",
        body: JSON.stringify({
          id: selectedItem.metadata.guid
        })
      });
      if (status === 200) {
        setTokenUpgraded(true);
      }
    },
    []
  );
  return (
    
    <Router history={history}>
      <Navbar
        isLoaded={!isLoadingAccounts}
        items={accounts}
        accountSelected={handleAccountChosen}
      />
      
        <Switch>
          <Route path="/" exact>
            <AppPage 
              hasChosenAccount={hasChosenAccount} 
              tokenUpgraded={tokenUpgraded} 
              accountID={accountID} 
            />
          </Route>
          <Route path="/login" exact component={Login}/>
          <Route path="/create" component={DumPage} />
        </Switch>
      </Router>
    
  );
};

export default AppRouter;
