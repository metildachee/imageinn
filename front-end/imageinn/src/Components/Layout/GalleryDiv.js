import React from 'react';

function GalleryDiv({ children, direction = "right", galleryItemsAspectRatio }) {
  return (
    <div className="gallery" data-direction={direction}>
      <div className="floating_content" data-images={galleryItemsAspectRatio}>
        {children}
      </div>
    </div>
  );
}

export default GalleryDiv;