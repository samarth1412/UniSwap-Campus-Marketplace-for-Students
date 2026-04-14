import { useCallback, useEffect, useId, useMemo, useRef, useState } from 'react';
import './ImageUpload.css';

const DEFAULT_MAX_BYTES = 5 * 1024 * 1024;

function isAllowedImage(file: File): boolean {
  if (file.type && file.type.startsWith('image/')) return true;
  const name = file.name.toLowerCase();
  return /\.(jpe?g|png|gif|webp|bmp|svg)$/i.test(name);
}

function formatBytes(size: number): string {
  if (size < 1024) return `${size} B`;
  const kb = size / 1024;
  if (kb < 1024) return `${kb.toFixed(1)} KB`;
  return `${(kb / 1024).toFixed(2)} MB`;
}

export interface ImageUploadProps {
  label?: string;
  file: File | null;
  onFileChange: (file: File | null) => void;
  existingImageUrl?: string | null;
  disabled?: boolean;
  disabledReason?: string;
  /** Max file size in bytes (default 5 MB). */
  maxSizeBytes?: number;
}

export function ImageUpload({
  label = 'Listing image',
  file,
  onFileChange,
  existingImageUrl = null,
  disabled = false,
  disabledReason = 'Image upload is temporarily disabled.',
  maxSizeBytes = DEFAULT_MAX_BYTES,
}: ImageUploadProps) {
  const inputId = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragOver, setDragOver] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

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
      setStatusMessage(null);
      if (!next) {
        onFileChange(null);
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      if (!isAllowedImage(next)) {
        setLocalError(`"${next.name}" is not a supported image. Please choose JPG, PNG, GIF, or WebP.`);
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      if (next.size > maxSizeBytes) {
        const mb = Math.round(maxSizeBytes / (1024 * 1024));
        setLocalError(
          `Image is too large (${formatBytes(next.size)}). Maximum allowed size is ${mb} MB.`
        );
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      onFileChange(next);
      setStatusMessage(`Selected "${next.name}" (${formatBytes(next.size)}).`);
    },
    [maxSizeBytes, onFileChange]
  );

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selected = event.target.files?.[0] ?? null;
    applyFile(selected);
  };

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault();
    if (disabled) return;
    setDragOver(false);
    const dropped = event.dataTransfer.files?.[0];
    if (dropped) applyFile(dropped);
  };

  const openPicker = () => {
    if (!disabled) inputRef.current?.click();
  };

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
        disabled={disabled}
        onChange={handleInputChange}
      />

      <div
        role="button"
        tabIndex={disabled ? -1 : 0}
        aria-disabled={disabled}
        className={`image-upload__drop${dragOver ? ' image-upload__drop--drag' : ''}${localError ? ' image-upload__drop--error' : ''}`}
        style={{ cursor: disabled ? 'not-allowed' : 'pointer', opacity: disabled ? 0.6 : 1 }}
        onClick={openPicker}
        onKeyDown={(e) => {
          if (disabled) return;
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            openPicker();
          }
        }}
        onDragEnter={(e) => {
          e.preventDefault();
          if (disabled) return;
          setDragOver(true);
        }}
        onDragOver={(e) => {
          e.preventDefault();
          if (disabled) return;
          e.dataTransfer.dropEffect = 'copy';
        }}
        onDragLeave={(e) => {
          if (disabled) return;
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
              disabled={disabled}
              onClick={(e) => {
                e.stopPropagation();
                applyFile(null);
              }}
            >
              {file ? 'Remove' : 'Choose another'}
            </button>
            <button
              type="button"
              className="image-upload__remove"
              disabled={disabled}
              onClick={(e) => {
                e.stopPropagation();
                openPicker();
              }}
            >
              Replace
            </button>
          </div>
        )}
      </div>

      {localError && <p className="image-upload__error">{localError}</p>}
      {!localError && statusMessage && (
        <p className="image-upload__hint" aria-live="polite">
          {statusMessage}
        </p>
      )}
      {disabled && (
        <p className="image-upload__hint" aria-live="polite">
          {disabledReason}
        </p>
      )}

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
