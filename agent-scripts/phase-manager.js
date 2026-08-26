#!/usr/bin/env node
import { runPhaseManager } from './git-helpers.js';
runPhaseManager(process.argv.slice(2));
