import React, { useState } from "react";
import { Checkbox } from "antd";

const CheckboxWithLabel = ({ label, onChange, disable }) => {
  const [disabled, setDisabled] = useState(false);
  return (
    <Checkbox
      style={{ color: "#5AA6FB" }}
      className="custom-checkbox"
      onChange={onChange}
    >
      <span
        className="checkbox-label"
        // style={{ borderColor: disabled ? "gray" : "black" }}
      >
        <span
          style={{ color: "#8A8FEA" }}
          className="cormorant-garamond-bold-italic"
        >
          {label}
        </span>
      </span>
    </Checkbox>
  );
};
export default CheckboxWithLabel;
