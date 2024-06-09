import React, { useState } from "react";
import { PlusOutlined, FileTextOutlined, FileImageOutlined } from "@ant-design/icons";
import { Button, Cascader, Checkbox, ColorPicker, DatePicker, Form, Input, InputNumber, Radio, Divider, Modal } from "antd";
import TokenInput from "./TokenInput";
import BoldTextSemiBold from "./BoldTextSemiBold";
import CheckboxWithLabel from "./Checkbox";
import IconComponent from "./Icon";

const { RangePicker } = DatePicker;
const { TextArea } = Input;

const normFile = (e) => {
  if (Array.isArray(e)) {
    return e;
  }
  return e?.fileList;
};

const SearchFilterForm = () => {
  const [useText, setUseText] = useState(true);
  const [imgDisabled, setImgDisabled] = useState(false);
  const [formValues, setFormValues] = useState({});
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleFormChange = (changedValues, allValues) => {
    setFormValues(allValues);
  };

  const handleSubmit = () => {
    console.log("Form Values: ", formValues);
    // You can propagate the form values here
    setIsModalOpen(false);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <Button type="primary" onClick={() => setIsModalOpen(true)}>
        Open Filters
      </Button>
      <Modal
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={handleCancel}
        width={500}
        style={{ borderColor: "black", borderWidth: 1, borderRadius: 5 }}
      >
        <BoldTextSemiBold bold="SEARCH" semiBold="FILTERS" fontSize="50px" />
        <Divider />
        <Radio.Group className="custom-radio">
          <Radio
            value="text"
            checked={useText}
            onChange={(e) => setUseText(e.target.checked)}
          >
            <FileTextOutlined style={{ fontSize: "25px" }} />
            <BoldTextSemiBold fontSize="30px" bold={"IN"} semiBold={"TEXT"} />
          </Radio>
          <Divider />
          <Form
            layout="horizontal"
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
              <CheckboxWithLabel disabled={useText} label="FUZZY" />
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
            >
              <TokenInput />
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

export default SearchFilterForm;
