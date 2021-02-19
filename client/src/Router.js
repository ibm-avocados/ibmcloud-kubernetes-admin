import React, { useCallback, useEffect, useState, Suspense } from 'react';
import { Loading } from 'carbon-components-react';
import Navbar from './common/Navbar';
import history from './globalHistory';
import queryString from 'query-string';

const AppPage = React.lazy(() => import('./pages/clusters/AppPage'));
const CreatePage = React.lazy(() => import('./pages/create/CreatePage'));
const SchedulePage = React.lazy(() => import('./pages/schedule/SchedulePage'));
const SettingsPage = React.lazy(() => import('./pages/settings/SettingsPage'));
const Login = React.lazy(() => import('./Login'));

const ApplicationRouter = (props) => {
  const query = queryString.parse(props.location.search);
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();
  const [selectedAccount, setSelectedAccount] = useState();
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
      const { location } = props;

      history.push(
        location.pathname + '?account=' + selectedItem.metadata.guid
      );
      history.go();
      // setSelectedAccount(selectedItem);
      // setAccountStuff(selectedItem.metadata.guid);
    },
    [setAccountStuff]
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
      if (query.account) {
        const item = accounts.resources.find(
          (account) => account.metadata.guid === query.account
        );
        if (item) {
          setSelectedAccount(item);
          setAccountStuff(item.metadata.guid);
        }
      }
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
        selectedItem={selectedAccount}
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
          query={query}
          hasChosenAccount={hasChosenAccount}
          tokenUpgraded={tokenUpgraded}
          accountID={accountID}
        />
      </Route>
    </>
  );
};

const App = ({location}) => {
  const {pathname, search} = location;
  
  useEffect(() => {
    fetch('/api/v1/login').then(({ status }) => {
      if (status !== 200) {
        if (pathname !== '/login'){
          history.push(`/login?state=${encodeURIComponent(pathname+search)}`)
        } else {
          history.push('/login');
        }
      }
    });
  }, []);
  //style={path === “create” ? styles.activeItem : styles.item}

  return (
    <Suspense fallback={<Loading />}>
      <Switch>
        <Route path="/login" exact component={Login}/>
        <Route path="/" component={ApplicationRouter} />
      </Switch>
    </Suspense>
  );
};

const AppRouter = () => {
  return (
    <Router history={history}>
      <Route path="/" component={App}/>
    </Router>
  );
};

export default AppRouter;
