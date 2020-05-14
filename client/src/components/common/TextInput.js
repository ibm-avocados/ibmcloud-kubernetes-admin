import React from 'react';
import styles from './TextInput.module.css';

const TextInput = (props) => (
  <input
    className={styles.input}
    value={props.value}
    onChange={props.onChange}
    placeholder={props.placeholder}
    onKeyDown={props.onKeyDown}
    style={props.style}
  />
);

export default TextInput;
