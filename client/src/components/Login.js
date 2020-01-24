import React from "react";
import Button from "./common/Button";
import styles from "./Login.module.css";

const Login = (props) => {
  const onClick = async () =>  {
    let response = await fetch("/api/v1/login");
    console.log(response.status);
    if (response.status > 400) {
      props.onResult(2);
      let endpointResponse = await fetch("/api/v1/identity-endpoints");
      let endpoints = await endpointResponse.json();

      window.open(endpoints.passcode_endpoint, "_blank");
      return;
    } else if(response.status === 200) {
      // session found and ok to log in
      console.log("logged in");
      props.getAccounts();
      props.onResult(1);
      return
    } else {
      props.onResult(0);
    }
  
  };
  return (
    <>
      <div className={styles.wrapper}>
        <Button label="Login with IBMId" onClickHandler={onClick} />
      </div>
    </>
  );
};

export default Login;
