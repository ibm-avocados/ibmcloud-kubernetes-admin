import React, { useCallback, useEffect, useState } from "react";
import { getJSON } from "./fetchUtil";
import AppPage from "./components/AppPage";
import LoginPage from "./components/LoginPage";
import OneTimePasscodePage from "./components/OneTimePasscodePage";

const STATE_INIT = 0;
const STATE_IS_LOGGED_IN = 1;
const STATE_SHOW_OTP = 2;

const App = () => {
  const [loginState, setLoginState] = useState(STATE_INIT);

  useEffect(() => {
    fetch("/api/v1/login").then(({ status }) => {
      if (status === 200) {
        setLoginState(STATE_IS_LOGGED_IN);
      }
    });
  }, []);

  const handleLoginClick = useCallback(async () => {
    setLoginState(STATE_SHOW_OTP);
    const endpoints = await fetch("/api/v1/identity-endpoints").then(getJSON);
    window.open(endpoints.passcode_endpoint, "_blank");
  }, []);

  const handleOTPSubmit = useCallback(async otp => {
    const { status } = await fetch("/api/v1/authenticate", {
      method: "POST",
      body: JSON.stringify({
        otp: otp
      })
    });

    if (status === 200) {
      setLoginState(STATE_IS_LOGGED_IN);
      return;
    }

    setLoginState(STATE_INIT);
  }, []);

  switch (loginState) {
    case STATE_SHOW_OTP:
      return <OneTimePasscodePage onSubmit={handleOTPSubmit} />;
    case STATE_IS_LOGGED_IN:
      return <AppPage />;
    default:
      return <LoginPage onLoginClick={handleLoginClick} />;
  }
};

export default App;
