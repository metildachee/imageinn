import React, { useState } from "react";
import { Checkbox } from "antd";

const CheckboxWithLabel = ({ label, onChange, disable }) => {
  const [disabled, setDisabled] = useState(false);
  return (
    <Checkbox
      style={{ color: "#f8dce4" }}
      className="custom-checkbox"
      onChange={onChange}
    >
      <span
        className="checkbox-label"
        // style={{ borderColor: disabled ? "gray" : "black" }}
      >
        <span
          style={{ color: "white" }}
          className="cormorant-garamond-bold-italic"
        >
          {label}
        </span>
      </span>
    </Checkbox>
  );
};
export default CheckboxWithLabel;
