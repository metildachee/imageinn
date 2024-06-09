import React from 'react';
import CategoryButton from './CategoryButton';

const CategoryButtonsRow = ({ categories }) => {
  console.log(typeof(categories))
  return (
    <div style={{ display: 'flex', justifyContent: 'space-around', flexWrap: 'wrap' }}>
      {categories.map((button, idx) => (
        <CategoryButton
          key_string={button.key_string}
          count={button.count}
        />
      ))}
    </div>
  );
};

export default CategoryButtonsRow;