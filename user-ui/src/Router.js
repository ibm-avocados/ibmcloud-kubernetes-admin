import React from 'react';
import { Router, Switch, Route } from 'react-router-dom';
import Home from './pages/Home/Home';

const AppRouter = () => {
    return(
        <Router>
            <Route path="/" component={Home} />
            <Router path="/:data" />
        </Router>
    )
}

export default AppRouter;