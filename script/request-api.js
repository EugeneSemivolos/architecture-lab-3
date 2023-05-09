const sendHTTPRequest = async (url, data) => {
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
      },
      body: data
    })
    if (response.ok) return response.text();
    throw new Error(`Request failed with status ${response.status}`);
  } catch (error) {
    console.error(error);
  }
};

export { sendHTTPRequest };
