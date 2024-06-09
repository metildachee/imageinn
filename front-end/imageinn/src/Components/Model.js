import React, { useState } from "react";
import { Button, Modal, Form, Radio, Divider } from "antd";
import { PlusOutlined, FileTextOutlined, FileImageOutlined } from "@ant-design/icons";
import MultiInputButton from "./MultiInputButton";
import CheckboxWithLabel from "./Checkbox";
import BoldTextSemiBold from "./BoldTextSemiBold";
import TokenInput from "./TokenInput";

const Model = ({ query, onFormSubmit }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [useText, setUseText] = useState(true);
  const [imgDisabled, setImgDisabled] = useState(false);
  const [tokens, setTokens] = useState([]);
  const [form] = Form.useForm();

  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleOk = () => {
    form.validateFields()
      .then(values => {
        const formData = { ...values, excludes: tokens };
        console.log('Form Values:', formData);
        setIsModalOpen(false);
        onFormSubmit(formData);
      })
      .catch(info => {
        console.log('Validate Failed:', info);
      });
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  const handleFormChange = (changedValues, allValues) => {
    // This can be used to update any state if needed
  };

  return (
    <>
      <MultiInputButton
        onClick={showModal}
        leftBold={"SEARCH"}
        leftSemibold={"FILTERS"}
        rightBold={query}
        backgroundColor={"white"}
      />
      <Modal
        open={isModalOpen}
        onOk={handleOk}
        onCancel={handleCancel}
        width={500}
        style={{
          border: "1px solid black",
          borderRadius: "15px",
          overflow: "hidden",
        }}
      >
        <BoldTextSemiBold bold="SEARCH" semiBold="FILTERS" fontSize="50px" />
        <Radio.Group className="custom-radio">
          <Radio
            value="text"
            checked={useText}
            onChange={(e) => setUseText(e.target.checked)}
          >
            <FileTextOutlined style={{ fontSize: "25px" }} />
            <BoldTextSemiBold fontSize="30px" bold={"IN"} semiBold={"TEXT"} />
          </Radio>

          <Form
            layout="horizontal"
            form={form}
            disabled={!useText}
            style={{ maxWidth: "100%" }}
            onValuesChange={handleFormChange}
          >
            <Form.Item
              style={{ padding: "0px", margin: "0px" }}
              label=""
              name="fuzzy"
              valuePropName="checked"
            >
              <CheckboxWithLabel label="FUZZY" />
            </Form.Item>

            <Form.Item
              className="roboto-medium"
              style={{ padding: "0px", margin: "0px" }}
              label=""
              name="is_and"
            >
              <Radio.Group>
                <Radio value="and"> AND </Radio>
                <Radio value="or"> OR </Radio>
              </Radio.Group>
            </Form.Item>
            <Form.Item
              style={{ padding: "0px", margin: "0px" }}
              className="roboto-medium"
              label="EX"
              name="excludes"
            >
              <TokenInput onTokensChange={setTokens} />
            </Form.Item>
          </Form>

          <Divider />
          <Radio
            value="img"
            checked={imgDisabled}
            onChange={(e) => {
              setImgDisabled(e.target.checked);
              setUseText(false);
            }}
          >
            <FileImageOutlined style={{ fontSize: "25px" }} />
            <BoldTextSemiBold fontSize="30px" bold={"IN"} semiBold={"IMAGE"} />
          </Radio>
        </Radio.Group>
      </Modal>
    </>
  );
};

export default Model;
