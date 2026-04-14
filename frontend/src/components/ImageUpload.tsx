import { useCallback, useEffect, useId, useMemo, useRef, useState } from 'react';
import './ImageUpload.css';

const DEFAULT_MAX_BYTES = 5 * 1024 * 1024;

function isAllowedImage(file: File): boolean {
  if (file.type && file.type.startsWith('image/')) return true;
  const name = file.name.toLowerCase();
  return /\.(jpe?g|png|gif|webp|bmp|svg)$/i.test(name);
}

export interface ImageUploadProps {
  label?: string;
  file: File | null;
  onFileChange: (file: File | null) => void;
  existingImageUrl?: string | null;
  /** Max file size in bytes (default 5 MB). */
  maxSizeBytes?: number;
}

export function ImageUpload({
  label = 'Listing image',
  file,
  onFileChange,
  existingImageUrl = null,
  maxSizeBytes = DEFAULT_MAX_BYTES,
}: ImageUploadProps) {
  const inputId = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragOver, setDragOver] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);

  const previewUrl = useMemo(() => {
    if (!file) return null;
    return URL.createObjectURL(file);
  }, [file]);
  const displayedImageUrl = previewUrl ?? existingImageUrl;

  useEffect(() => {
    return () => {
      if (previewUrl) URL.revokeObjectURL(previewUrl);
    };
  }, [previewUrl]);

  const applyFile = useCallback(
    (next: File | null) => {
      setLocalError(null);
      if (!next) {
        onFileChange(null);
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      if (!isAllowedImage(next)) {
        setLocalError('Please choose an image file (e.g. JPG, PNG, GIF, or WebP).');
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      if (next.size > maxSizeBytes) {
        const mb = Math.round(maxSizeBytes / (1024 * 1024));
        setLocalError(`Image must be ${mb} MB or smaller.`);
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      onFileChange(next);
    },
    [maxSizeBytes, onFileChange]
  );

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selected = event.target.files?.[0] ?? null;
    applyFile(selected);
  };

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault();
    setDragOver(false);
    const dropped = event.dataTransfer.files?.[0];
    if (dropped) applyFile(dropped);
  };

  const openPicker = () => inputRef.current?.click();

  return (
    <div className="image-upload">
      <span className="image-upload__label" id={`${inputId}-label`}>
        {label}
      </span>
      <p className="image-upload__hint">
        Optional — JPG, PNG, or WebP. Max {Math.round(maxSizeBytes / (1024 * 1024))} MB.
      </p>

      <input
        ref={inputRef}
        id={inputId}
        className="image-upload__input"
        type="file"
        accept="image/*"
        aria-labelledby={`${inputId}-label`}
        onChange={handleInputChange}
      />

      <div
        role="button"
        tabIndex={0}
        className={`image-upload__drop${dragOver ? ' image-upload__drop--drag' : ''}${localError ? ' image-upload__drop--error' : ''}`}
        onClick={openPicker}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            openPicker();
          }
        }}
        onDragEnter={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragOver={(e) => {
          e.preventDefault();
          e.dataTransfer.dropEffect = 'copy';
        }}
        onDragLeave={(e) => {
          if (!e.currentTarget.contains(e.relatedTarget as Node)) setDragOver(false);
        }}
        onDrop={handleDrop}
      >
        {!file && !existingImageUrl ? (
          <>
            <p className="image-upload__cta">Drop an image here or click to browse</p>
            <p className="image-upload__cta-secondary">Your photo helps buyers trust the listing.</p>
          </>
        ) : (
          <div className="image-upload__file-row">
            {file ? (
              <p className="image-upload__file-name" title={file.name}>
                Selected: {file.name}
              </p>
            ) : (
              <p className="image-upload__file-name">Current image selected. Choose a new one to replace it.</p>
            )}
            <button
              type="button"
              className="image-upload__remove"
              onClick={(e) => {
                e.stopPropagation();
                applyFile(null);
              }}
            >
              {file ? 'Remove' : 'Choose another'}
            </button>
          </div>
        )}
      </div>

      {localError && <p className="image-upload__error">{localError}</p>}

      {displayedImageUrl && (
        <div className="image-upload__preview-wrap">
          <img
            className="image-upload__preview"
            src={displayedImageUrl}
            alt={file ? 'Preview of selected listing image' : 'Current listing image'}
          />
        </div>
      )}
    </div>
  );
}
