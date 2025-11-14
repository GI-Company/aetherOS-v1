import React, { useState } from 'react';
import { getAuth, sendSignInLinkToEmail } from 'firebase/auth';

function Login() {
  const [email, setEmail] = useState('');
  const [isLinkSent, setIsLinkSent] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const auth = getAuth();
    const actionCodeSettings = {
      // URL you want to redirect back to. The domain (www.example.com) must be
      // authorized in the Firebase console.
      url: window.location.href, // Using current URL for simplicity
      handleCodeInApp: true, // This must be true.
    };

    try {
      await sendSignInLinkToEmail(auth, email, actionCodeSettings);
      // The link was successfully sent. Inform the user.
      // Save the email locally so you don't need to ask the user for it again
      // if they open the link on the same device.
      window.localStorage.setItem('emailForSignIn', email);
      setIsLinkSent(true);
    } catch (error) {
      setError(error.message);
    }
  };

  if (isLinkSent) {
    return (
      <div>
        <h2>Link Sent!</h2>
        <p>A sign-in link has been sent to {email}. Please check your inbox.</p>
      </div>
    );
  }

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Enter your email"
          required
        />
        <button type="submit">Send Sign-In Link</button>
      </form>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
}

export default Login;
