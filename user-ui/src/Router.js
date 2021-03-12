import React, {useState, useEffect} from 'react';
import { Router, Route, Switch } from 'react-router-dom';
import Home from './pages/Home/Home';
import Workshop from './pages/Workshop/Workshop';
import history from './globalHistory';
import Header from './components/Header/Header';

const App = ({location}) => {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState({});
  useEffect(() => {
    fetch('/api/v1/login').then(({status})=> {
      if(status !== 200) {
        setLoggedIn(false);
      } else {
        setLoggedIn(true);
      }
    });
  }, []);

  //TODO: Get the api for setting user info
  return (
    <>
      <Header location={location} loggedIn={loggedIn} userName="Mofi Rahman" />
      <Switch>
        <Route path="/:workshop" component={Workshop} />
        <Route path="/" component={Home} />
      </Switch>
    </>
  );
};

const AppRouter = () => {
  return (
    <Router history={history}>
      <Route path="/" component={App} />
    </Router>
  );
};

export default AppRouter;