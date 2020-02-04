import React, {useState, useCallback} from 'react';
import Button from "./common/Button";
import TextInput from "./common/TextInput";
import styles from "./pagestyles.module.css";

const LoginPage = ({ onLoginClick }) => (
  <div className={styles.wrapper}>
    <Button label="Login with IBMId" onClickHandler={onLoginClick} />
  </div>
);

const OneTimePasscodePage = ({ onSubmit }) => {
  const [value, setValue] = useState("");

  const handleChange = useCallback(e => {
    setValue(e.target.value);
  }, []);

  const handleKeyDown = useCallback(
    e => {
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


const Login = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [attemtedLogin, setAttemtedLogin] = useState(false);

  useEffect(() => {
    fetch("/api/v1/login").then(({ status }) => {
      if (status === 200) {
        setIsLoggedIn(true);
      }
    });
  }, []);

  const handleLoginClick = useCallback(async () => {
    setAttemtedLogin(true);
    const response = await fetch("/api/v1/identity-endpoints");
    const endpoints = await response.json();
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
      setIsLoggedIn(true);
      return;
    }

    setIsLoggedIn(false);
  }, []);

  if (isLoggedIn) {
    history.push("/");
  } else if (attemtedLogin) {
    return (<OneTimePasscodePage onSubmit={handleOTPSubmit}/>)
  }
  return (<LoginPage onLoginClick={handleLoginClick}/>)
};

export default Login;




