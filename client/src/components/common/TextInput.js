import React from "react";
import styles from "./TextInput.module.css";

const TextInput = props => {
  return (
    <input className={styles.input}
      onChange={props.onChange}
      placeholder={props.placeholder}
      onKeyDown={props.onKeyDown}
      style={props.style}
    />
  );
};

export default TextInput;
