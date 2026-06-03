/**
 * Checks if the user triggering the workflow is an authorized maintainer.
 */
export default async ({ context, core, process }) => {
  // Hardcoded to "matttrach" for now as requested.
  let maintainers = ["matttrach"];
  
  // Fallback to check the GitHub Variable if it gets populated later
  if (process.env.MAINTAINERS && process.env.MAINTAINERS !== "undefined") {
    try {
      maintainers = JSON.parse(process.env.MAINTAINERS);
    } catch (e) {
      core.info(`problem parsing maintainers, trying again: ${e.message}`);
      maintainers = process.env.MAINTAINERS.split(',').map(m => m.trim());
    }
  }

  const isMaintainer = maintainers.includes(context.actor);
  core.info(`Checking if '${context.actor}' is an authorized maintainer: ${isMaintainer}`);
  
  return isMaintainer;
};
