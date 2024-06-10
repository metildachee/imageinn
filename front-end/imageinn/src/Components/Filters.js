import React, { useState, useEffect } from "react";
import { Radio, Form, Space, Row, Col } from "antd";
import { AlignRightOutlined } from "@ant-design/icons";
import TokenInput from "./TokenInput";
import CheckboxWithLabel from "./Checkbox";

const RectangleComponent = ({
  updateFuzzy,
  updateIsAnd,
  updateExcludes,
  updateImage,
  updatePersonalise
}) => {
  const [useText, setUseText] = useState(true);
  const [imgDisabled, setImgDisabled] = useState(false);
  const [tokens, setTokens] = useState([]);
  const [form] = Form.useForm();
  const [isOr, setIsOr] = useState(true);
  const [fuzzy, setFuzzy] = useState(false);
  const [isImg, setImg] = useState(false);
  const [isPersonalise, setPersonalise] = useState(false);

  useEffect(() => {
    updateFuzzy(fuzzy);
    updateIsAnd(!isOr);
    updateExcludes(tokens);
    updateImage(isImg)
  }, [form.getFieldsValue()]);

  return (
    <Space align="center">
      <div
        className="rectangle custom-radio"
        style={{
          border: "1px solid #8A8FEA",
          borderRadius: "5px",
          color: "#8A8FEA",
          textAlign: "left",
          padding: "10px",
          maxWidth: "100%",
          marginLeft: "10px",
          position: "relative", // Add position relative for positioning the icon
          backgroundColor: "white"
        }}
      >
        {/* Icon for aligning right */}
        <div
          style={{
            position: "absolute",
            top: "-20px",
            right: "5px",
          }}
        >
          <AlignRightOutlined />
        </div>

        <Form
          layout="horizontal"
          form={form}
          disabled={imgDisabled}
          style={{ width: "500px" }}
        >
          <Row>
            <Col span={6}>
              <Form.Item label="" style={{ margin: 0 }}>
                <CheckboxWithLabel
                  label="FUZZY"
                  disabled={imgDisabled}
                  onChange={(e) => setFuzzy(e.target.checked)}
                />
              </Form.Item>
            </Col>
            <Col span={10}>
              <Form.Item
                label=""
                style={{ margin: 0, color: "#8A8FEA" }}
                className="cormorant-garamond-regular-italic"
              >
                <Radio.Group>
                  <Radio
                    className="custom-radio"
                    value="and"
                    style={{ color: "#8A8FEA" }}
                    disabled={imgDisabled}
                    onChange={(e) => {
                      console.log(e.target.value);
                      setIsOr(false);
                      updateIsAnd(false);
                    }}
                  >
                    {" "}
                    AND{" "}
                  </Radio>
                  <Radio
                    className="custom-radio"
                    value="or"
                    style={{ color: "#8A8FEA" }}
                    disabled={imgDisabled}
                    onChange={(e) => {
                      console.log(e.target.value);
                      setIsOr(true);
                      updateIsAnd(true);
                    }}
                  >
                    {" "}
                    OR{" "}
                  </Radio>
                </Radio.Group>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item label="" style={{ margin: 0 }}>
                <CheckboxWithLabel
                  label={"IN IMAGE"}
                  onChange={(e) => {
                    // setImgDisabled(e.target.checked);
                    // setUseText(!e.target.checked);
                    console.log(e.target.checked)
                    updateImage(e.target.checked);
                    setImg(e.target.checked)
                  }}
                />
              </Form.Item>
              <Form.Item label="" style={{ margin: 0 }}>
                <CheckboxWithLabel
                  label={"PERSONALISE"}
                  onChange={(e) => {
                    // setImgDisabled(e.target.checked);
                    // setUseText(!e.target.checked);
                    updatePersonalise(e.target.checked);
                    setPersonalise(e.target.checked)
                  }}
                />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <Form.Item label="" style={{ margin: 0 }}>
                <TokenInput onTokensChange={setTokens} disabled={imgDisabled} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </div>
      <div
        style={{
          position: "absolute",
          bottom: -3,
          right: 5,
          width: `45ch`,
          height: `70px`,
          borderRight: "1px solid #8A8FEA",
          borderBottom: "1px solid #8A8FEA",
          borderTopColor: "#DC648F",
          borderLeftColor: "#DC648F",
          boxSizing: "border-box",
          overflow: "hidden", // Hide the top and left borders
          pointerEvents: "none",
          margin: "-2px",
          borderBottomRightRadius: "5px",
        }}
      />
    </Space>
  );
};

export default RectangleComponent;
