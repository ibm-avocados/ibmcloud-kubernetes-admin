import React from 'react';
import {Router, Route, Switch} from 'react-router-dom';
import Home from './pages/Home/Home';
import Workshop from './pages/Workshop/Workshop';
import history from './globalHistory';
import Header from './components/Header/Header';

const App = () => {
    return (
        <>
            <Header loggedIn={true}/>
            <Switch>
                <Route path="/:workshop" component={Workshop}/>
                <Route path="/" component={Home}/>
            </Switch>
        </>
    )
}

const AppRouter = () => {
    return (
        <Router history={history}>
            <App/>
        </Router>
    );
}

export default AppRouter;