import React from 'react';
import { Button, Flex } from 'antd';
const CategoryButton = ({ key_string, count }) => (
  <Flex gap="small" wrap>
    <Button>{key_string} {count}</Button>
  </Flex>
);
export default CategoryButton;