import React, { useState } from "react";
import { Tag, Input } from "antd";
import { PlusOutlined } from "@ant-design/icons";

const TokenInput = ({ onTokensChange }) => {
  const [inputValue, setInputValue] = useState("");
  const [tokens, setTokens] = useState([]);

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleInputConfirm = () => {
    if (inputValue && tokens.length < 3) {
      const newTokens = [...tokens, inputValue];
      setTokens(newTokens);
      setInputValue("");
      onTokensChange(newTokens);
    }
  };

  const handleTokenClose = (removedToken) => {
    const newTokens = tokens.filter((token) => token !== removedToken);
    setTokens(newTokens);
    onTokensChange(newTokens);
  };

  return (
    <div>
      <div style={{ marginBottom: 8, height: "20px" }}>
        {tokens.map((token, index) => (
          <Tag
            style={{ color: "#8A8FEA", borderColor: "#8A8FEA", marginRight: 4 }}
            className="cormorant-garamond-bold-italic"
            key={index}
            closable
            onClose={() => handleTokenClose(token)}
          >
            {token}
          </Tag>
        ))}
      </div>
      <Input
        prefix="EXCLUDE"
        style={{
          color: "#white",
          border: "1px solid #f8dce4",
          backgroundColor: "transparent",
        //   marginTop: 24, // Add padding to the top
        }}
        value={inputValue}
        onChange={handleInputChange}
        onPressEnter={handleInputConfirm}
        className="cormorant-garamond-bold-italic input-with-placeholder"
        // placeholder={`${3 - tokens.length} keywords`}
      />
    </div>
  );
};

export default TokenInput;
