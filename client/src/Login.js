import React, { useState, useCallback, useEffect } from "react";
import Button from "./common/Button";
import TextInput from "./common/TextInput";
import styles from "./pagestyles.module.css";
import history from "./globalHistory";
import queryString from "query-string";

const OneTimePasscodePage = ({ onSubmit }) => {
  const [value, setValue] = useState("");

  const handleChange = useCallback((e) => {
    setValue(e.target.value);
  }, []);

  const handleKeyDown = useCallback(
    (e) => {
      if (e.key === "Enter") {
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

const Login = ({ location }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [attemtedLogin, setAttemtedLogin] = useState(false);
  const [identityEndpoint, setIdentityEndpoint] = useState(undefined);

  useEffect(() => {
    fetch("/api/v1/login").then(({ status }) => {
      if (status === 200) {
        setIsLoggedIn(true);
      }
    });
    fetch("/api/v1/identity-endpoints")
      .then((r) => r.json())
      .then(({ passcode_endpoint }) => {
        setIdentityEndpoint(passcode_endpoint);
      });
  }, []);

  const handleOTPSubmit = useCallback(async (otp) => {
    const { status } = await fetch("/api/v1/authenticate", {
      method: "POST",
      body: JSON.stringify({
        otp,
      }),
      headers: {
        'Content-Type': 'application/json'
      },
    });

    if (status === 200) {
      setIsLoggedIn(true);
      return;
    }

    setIsLoggedIn(false);
  }, []);

  if (isLoggedIn) {
    // const queryData = queryString.stringify(query);
    const { search } = location;
    const query = queryString.parse(search);

    console.log(query);

    history.push(query.state);
  } else if (attemtedLogin) {
    return <OneTimePasscodePage onSubmit={handleOTPSubmit} />;
  }
  return (
    <div className={styles.wrapper}>
      {identityEndpoint && (
        <div className={styles.center}>
          <a
            className={styles.button}
            href={identityEndpoint}
            onClick={() => setAttemtedLogin(true)}
            target="_blank"
            rel="noopener noreferrer"
          >
            Login with IBMId
          </a>
          <div className={styles.message}>Come back to this page with your OTP</div>
        </div>
      )}
    </div>
  );
};

export default Login;
