import React from "react";
import Button from "./common/Button";

import styles from "./Login.module.css";

const LoginPage = ({ onLoginClick }) => (
  <div className={styles.wrapper}>
    <Button label="Login with IBMId" onClickHandler={onLoginClick} />
  </div>
);

export default LoginPage;
