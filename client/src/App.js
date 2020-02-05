import React, { useEffect } from "react";
import Router from "./Router";
import history from "./globalHistory";

const App = () => {
  console.log("APP Invoked");
  useEffect(() => {
    fetch("/api/v1/login").then(({ status }) => {
      if (status !== 200) {
        history.push("/login");
      }
    });
  }, []);

  return <Router />
};

export default App;
