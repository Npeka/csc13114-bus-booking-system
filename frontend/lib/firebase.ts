import { initializeApp, getApps, FirebaseApp } from "firebase/app";
import {
  getAuth,
  Auth,
  RecaptchaVerifier,
  ApplicationVerifier,
} from "firebase/auth";

// Firebase configuration from environment variables
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
};

// Initialize Firebase app (singleton pattern)
let app: FirebaseApp | undefined;
let auth: Auth | undefined;

if (typeof window !== "undefined") {
  // Client-side only initialization
  if (!getApps().length) {
    app = initializeApp(firebaseConfig);
  } else {
    app = getApps()[0];
  }
  auth = getAuth(app);
}

// Recaptcha verifier instance (for phone auth)
let recaptchaVerifier: ApplicationVerifier | null = null;

/**
 * Get or create a RecaptchaVerifier instance for phone authentication
 * @param containerId - The ID of the DOM element to render the reCAPTCHA
 * @param isInvisible - Whether to use invisible reCAPTCHA
 */
export const getRecaptchaVerifier = (
  containerId: string = "recaptcha-container",
  isInvisible: boolean = false,
): ApplicationVerifier => {
  if (typeof window === "undefined") {
    throw new Error("RecaptchaVerifier can only be used in browser");
  }

  if (!auth) {
    throw new Error("Firebase auth not initialized");
  }

  // Check if container exists
  const container = document.getElementById(containerId);
  if (!container) {
    throw new Error(
      `reCAPTCHA container with id '${containerId}' not found in DOM`,
    );
  }

  // Return existing verifier if already created and valid
  if (recaptchaVerifier) {
    return recaptchaVerifier;
  }

  // Create new verifier
  recaptchaVerifier = new RecaptchaVerifier(auth, containerId, {
    size: isInvisible ? "invisible" : "normal",
    callback: () => {
      // reCAPTCHA solved - allow user to proceed
      console.log("reCAPTCHA verified");
    },
    "expired-callback": () => {
      // reCAPTCHA expired - reset
      console.log("reCAPTCHA expired");
      recaptchaVerifier = null;
    },
  });

  return recaptchaVerifier;
};

/**
 * Clear the RecaptchaVerifier instance
 */
export const clearRecaptchaVerifier = () => {
  recaptchaVerifier = null;
};

// Export auth instance
export { auth };
export default app;
