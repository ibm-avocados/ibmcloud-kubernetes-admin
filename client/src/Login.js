import React, { useState, useCallback, useEffect } from 'react';
import styles from './pagestyles.module.css';
import history from './globalHistory';
import queryString from 'query-string';

const Login = ({ location }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    fetch('/api/v1/login').then(({ status }) => {
      if (status === 200) {
        setIsLoggedIn(true);
      }
    });
  }, []);

  if (isLoggedIn) {
    // const queryData = queryString.stringify(query);
    const { search } = location;
    const query = queryString.parse(search);

    console.log(query);
    if (query.state === null || query.state === undefined) {
      history.push('/')
    }
    else {
      history.push(query.state);
    }
  }

  const getQuery = () => {
    const { search } = location;
    console.log('getquery', search);
    const query = queryString.parse(search);
    console.log('getquery', query.state);
    let data = parseQuery(query.state);
    console.log('getquery', data);
    return '/auth?provider=ibm&login=true&' + data;
  }

  const parseQuery = (data) => {
    if (data === null || data === undefined) {
      return '';
    }
    return data.split('/?')[1]
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.center}>
        <a
          className={styles.button}
          href={getQuery()}
        >
          Login with IBMId
        </a>
      </div>
    </div>
  );
};

export default Login;
