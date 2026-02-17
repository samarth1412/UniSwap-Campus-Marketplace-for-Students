import { useEffect, useState } from 'react';

interface ImageUploadProps {
  label?: string;
  file: File | null;
  onFileChange: (file: File | null) => void;
}

export function ImageUpload({ label = 'Listing image', file, onFileChange }: ImageUploadProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      return;
    }

    const url = URL.createObjectURL(file);
    setPreviewUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [file]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0] ?? null;
    onFileChange(selectedFile);
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
      <label style={{ fontWeight: 500 }}>
        {label}
        <input
          type="file"
          accept="image/*"
          onChange={handleChange}
          style={{ display: 'block', marginTop: '0.25rem' }}
        />
      </label>
      {previewUrl && (
        <img
          src={previewUrl}
          alt="Listing preview"
          style={{ maxWidth: '320px', borderRadius: '0.5rem', border: '1px solid #ddd' }}
        />
      )}
    </div>
  );
}

