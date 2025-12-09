import { create } from "zustand";

type DialogType = "login" | "register" | "phone" | "forgot-password" | null;

interface AuthDialogState {
  openDialog: DialogType;
  setOpenDialog: (dialog: DialogType) => void;
  closeAll: () => void;
}

export const useAuthDialog = create<AuthDialogState>((set) => ({
  openDialog: null,
  setOpenDialog: (dialog) => set({ openDialog: dialog }),
  closeAll: () => set({ openDialog: null }),
}));
