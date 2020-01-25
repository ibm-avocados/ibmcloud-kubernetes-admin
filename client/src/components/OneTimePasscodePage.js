import React, { useCallback, useState } from "react";
import TextInput from "./common/TextInput";

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
    <TextInput
      style={{
        position: "absolute",
        left: "50%",
        top: "50%",
        transform: "translate(-50%, -50%)"
      }}
      value={value}
      onChange={handleChange}
      placeholder="One Time Passcode"
      onKeyDown={handleKeyDown}
    />
  );
};

export default OneTimePasscodePage;
