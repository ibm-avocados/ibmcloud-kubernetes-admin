import React from 'react';
import {Router, Route, Switch} from 'react-router-dom';
import Home from './pages/Home/Home';
import Workshop from './pages/Workshop/Workshop';
import history from './globalHistory';

const AppRouter = () => {
    return (
        <Router history={history}>
            <Switch>
                <Route path="/:workshop" component={Workshop}/>
                <Route path="/" component={Home}/>
            </Switch>
        </Router>
    )
}

export default AppRouter;