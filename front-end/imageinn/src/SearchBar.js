import React from 'react'; // TODO: Change into a nicer category search bar
import { AudioOutlined } from '@ant-design/icons';
import { Input, Space } from 'antd';
const { Search } = Input;
const suffix = (
  <AudioOutlined
    style={{
      fontSize: 16,
      color: '#353741',
    }}
  />
);

const SearchBar = ({ onSearch }) => (
  <Space direction="vertical">
    <Search placeholder="input search text" onSearch={value => onSearch(value)} enterButton />
  </Space>
);
export default SearchBar;