CREATE TABLE IF NOT EXISTS users (
  user_id SERIAL PRIMARY KEY,
  email TEXT,
  name TEXT,
  github_id INTEGER UNIQUE NOT NULL,
  github_login TEXT,
  github_avatar_url TEXT,
  github_html_url TEXT,
  github_access_token TEXT,
  github_scope TEXT
);

--

CREATE OR REPLACE FUNCTION upsert_user(
  _email TEXT,
  _name TEXT,
  _github_id INTEGER,
  _github_login TEXT,
  _github_avatar_url TEXT,
  _github_html_url TEXT,
  _github_access_token TEXT,
  _github_scope TEXT
)
RETURNS integer
AS $$
DECLARE _user_id integer;
BEGIN
  -- It is practically impossible for a new user to request a new account with
  -- two simultaneous requests with different data, making this an acceptable
  -- method for upsert.
  UPDATE users
    SET
      email = _email,
      name = _name,
      github_login = _github_login,
      github_avatar_url = _github_avatar_url,
      github_html_url = _github_html_url,
      github_access_token = _github_access_token,
      github_scope = _github_scope
    WHERE github_id = _github_id;

  INSERT INTO users
    (
      email,
      name,
      github_id,
      github_login,
      github_avatar_url,
      github_html_url,
      github_access_token,
      github_scope
    )
    SELECT
      _email,
      _name,
      _github_id,
      _github_login,
      _github_avatar_url,
      _github_html_url,
      _github_access_token,
      _github_scope
    WHERE NOT EXISTS (SELECT 1 FROM users WHERE github_id = _github_id LIMIT 1);

  SELECT user_id INTO _user_id FROM users WHERE github_id = _github_id LIMIT 1;
  RETURN _user_id;
END;
$$ LANGUAGE plpgsql;
