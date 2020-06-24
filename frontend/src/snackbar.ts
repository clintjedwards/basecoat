export enum SnackBarColor {
  Error = "error",
  Success = "success",
}

export interface SnackBar {
  text: string;
  display: boolean;
  color: SnackBarColor;
}
