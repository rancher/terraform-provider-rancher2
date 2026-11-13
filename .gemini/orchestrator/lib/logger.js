export const log = (msg) => {
  console.log(`[INFO] ${new Date().toISOString()} - ${msg}`);
};

export const logError = (msg) => {
  console.error(`[ERROR] ${new Date().toISOString()} - ${msg}`);
};
