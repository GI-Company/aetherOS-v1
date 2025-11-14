import React, { useEffect, useState } from 'react';
import { getAuth, isSignInWithEmailLink, signInWithEmailLink } from 'firebase/auth';

function HandleAuth() {
  const [error, setError] = useState(null);
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const auth = getAuth();

    const completeSignIn = async () => {
      if (isSignInWithEmailLink(auth, window.location.href)) {
        let email = window.localStorage.getItem('emailForSignIn');
        if (!email) {
          // User opened the link on a different device. To prevent session
          // fixation attacks, ask the user to provide their email again.
          email = window.prompt('Please provide your email for confirmation');
        }

        try {
          const result = await signInWithEmailLink(auth, email, window.location.href);
          // Clear email from storage.
          window.localStorage.removeItem('emailForSignIn');
          
          // You can access the new user via result.user
          // Additional logic like sending the ID token to your backend would go here.
          const user = result.user;
          const idToken = await user.getIdToken();

          // TODO: Send idToken to your Aether backend to get a session token
          console.log("Successfully signed in!", user);
          console.log("Firebase ID Token:", idToken);

          // For now, we'll just redirect to the home page.
          window.location.href = "/";

        } catch (error) {
          setError(error.message);
        }
      } else {
        // This is not a sign-in link, so we don't need to do anything.
        setIsProcessing(false);
      }
    };

    completeSignIn();
  }, []);

  if (isProcessing && !error) {
    return <div><p>Completing sign-in...</p></div>;
  }

  if (error) {
    return (
      <div>
        <h2>Authentication Error</h2>
        <p>{error}</p>
        <p>Please try signing in again.</p>
      </div>
    );
  }

  // If it's not a sign-in link, we can just render nothing or a redirect.
  return null;
}

export default HandleAuth;
