import React, { useCallback, useEffect, useState } from 'react';
import { Router, Switch, Route } from 'react-router-dom';
import AppPage from './pages/clusters/AppPage';
import CreatePage from './pages/create/CreatePage';
import Login from './Login';
import Navbar from './common/Navbar';
import history from './globalHistory';

const HolderThing = () => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();

  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [tokenUpgraded, setTokenUpgraded] = useState(false);

  const setAccountStuff = useCallback(async (guid) => {
    localStorage.setItem('accountID', guid);
    setSelectedAccountID(guid);
    setTokenUpgraded(false);
    setHasChosenAccount(true);
    const { status } = await fetch('/api/v1/authenticate/account', {
      method: 'POST',
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
    [setAccountStuff],
  );

  useEffect(() => {
    const loadAccounts = async () => {
      setLoadingAccounts(true);
      const response = await fetch('/api/v1/accounts');
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
        isLoaded={!isLoadingAccounts}
        items={accounts}
        accountSelected={handleAccountChosen}
      />
      <Route path="/create" exact>
        <CreatePage tokenUpgraded={tokenUpgraded} hasChosenAccount={hasChosenAccount} accountID={accountID} />
      </Route>
      <Route path="/schedule" exact />
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
    fetch('/api/v1/login').then(({ status }) => {
      if (status !== 200) {
        history.push('/login');
      }
    });
  }, []);

  return (
    <Router history={history}>
      <Switch>
        <Route path="/login" exact component={Login} />
        <Route path="/">
          <HolderThing />
        </Route>
      </Switch>
    </Router>
  );
};

export default AppRouter;
