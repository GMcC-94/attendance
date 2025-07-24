export interface Image {
  id: number;
  fileName: string;
  filePath: string;
  type: string;
  uploadedAt: string;
}

export interface LogoUploadResponse {
  url: string;
}