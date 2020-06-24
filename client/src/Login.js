import React, { useState, useCallback, useEffect } from 'react';
import Button from './common/Button';
import TextInput from './common/TextInput';
import styles from './pagestyles.module.css';
import history from './globalHistory';

const OneTimePasscodePage = ({ onSubmit }) => {
  const [value, setValue] = useState('');

  const handleChange = useCallback((e) => {
    setValue(e.target.value);
  }, []);

  const handleKeyDown = useCallback(
    (e) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        e.stopPropagation();
        onSubmit(value);
      }
    },
    [onSubmit, value]
  );

  return (
    <div className={styles.wrapper}>
      <TextInput
        value={value}
        onChange={handleChange}
        placeholder="One Time Passcode"
        onKeyDown={handleKeyDown}
      />
    </div>
  );
};

const Login = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [attemtedLogin, setAttemtedLogin] = useState(false);
  const [identityEndpoint, setIdentityEndpoint] = useState(undefined);

  useEffect(() => {
    fetch('/api/v1/login').then(({ status }) => {
      if (status === 200) {
        setIsLoggedIn(true);
      }
    });
    fetch('/api/v1/identity-endpoints')
      .then((r) => r.json())
      .then(({ passcode_endpoint }) => {
        console.log(passcode_endpoint);
        setIdentityEndpoint(passcode_endpoint);
      });
  }, []);

  const handleOTPSubmit = useCallback(async (otp) => {
    const { status } = await fetch('/api/v1/authenticate', {
      method: 'POST',
      body: JSON.stringify({
        otp,
      }),
    });

    if (status === 200) {
      setIsLoggedIn(true);
      return;
    }

    setIsLoggedIn(false);
  }, []);

  if (isLoggedIn) {
    history.push('/');
  } else if (attemtedLogin) {
    return <OneTimePasscodePage onSubmit={handleOTPSubmit} />;
  }
  return (
    <div className={styles.wrapper}>
      {identityEndpoint && (
        <a href={identityEndpoint} onClick={() => setAttemtedLogin(true)} target="_blank" rel="noopener noreferrer">
          Login with IBMId
        </a>
      )}
    </div>
  );
};

export default Login;
