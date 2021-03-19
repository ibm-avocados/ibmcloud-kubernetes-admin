import React, { useState, useEffect } from 'react';
import { Router, Route, Switch } from 'react-router-dom';
import Home from './pages/Home/Home';
import Workshop from './pages/Workshop/Workshop';
import history from './globalHistory';
import Header from './components/Header/Header';

const App = ({ location }) => {
  const [loggedIn, setLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  useEffect(() => {
    fetch('/api/v1/login').then(({ status }) => {
      if (status !== 200) {
        setLoggedIn(false);
      } else {
        setLoggedIn(true);
      }
    });
  }, []);

  useEffect(() => {
    const getUserInfo = async () => {
      try {
        const response = await fetch('/api/v1/user/info');
        const data = await response.json();
        setUser(data);
      } catch (e) {
        console.log(e);
      }
    };

    if (loggedIn) {
      getUserInfo();
    }
  }, [loggedIn]);

  //TODO: Get the api for setting user info
  return (
    <>
      <Header location={location} loggedIn={loggedIn} user={user}/>
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