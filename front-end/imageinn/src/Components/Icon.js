import React from "react";
import Icon, { HomeOutlined } from '@ant-design/icons';
import { Space } from 'antd';

const IconComponent = ({ path, alt, width, height, fontSize }) => (
  <Icon
    component={() => (
      <img
        src={path}
        alt={alt}
        style={{ width: width, height: height, fontSize: fontSize }}
      />
    )}
  />
);
export default IconComponent;
