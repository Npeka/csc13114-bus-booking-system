"use client";

import { LoginDialog } from "./login-dialog";
import { RegisterDialog } from "./register-dialog";
import { PhoneLoginDialog } from "./phone-login-dialog";
import { ForgotPasswordDialog } from "./forgot-password-dialog";

interface AuthDialogsProps {
  loginOpen: boolean;
  onLoginOpenChange: (open: boolean) => void;
}

/**
 * Container component that renders all authentication dialogs.
 * Dialogs are managed through the useAuthDialog hook for seamless transitions.
 */
export function AuthDialogs({
  loginOpen,
  onLoginOpenChange,
}: AuthDialogsProps) {
  return (
    <>
      <LoginDialog isOpen={loginOpen} onOpenChange={onLoginOpenChange} />
      <RegisterDialog />
      <PhoneLoginDialog />
      <ForgotPasswordDialog />
    </>
  );
}
