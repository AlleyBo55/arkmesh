import { configureStore, createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type TerminalTone = "system" | "success" | "warning" | "muted";

export type TerminalEntry = {
  id: string;
  command?: string;
  lines: string[];
  tone: TerminalTone;
};

type TerminalState = {
  entries: TerminalEntry[];
  history: string[];
};

const initialState: TerminalState = {
  entries: [
    {
      id: "boot-0",
      lines: [
        "ARKMESH RECOVERY CONSOLE v0.1.0-alpha.1",
        "LOCAL EVIDENCE CHANNEL READY",
        "Type HELP or select a command. No capability is simulated as complete.",
      ],
      tone: "system",
    },
  ],
  history: [],
};

const terminalSlice = createSlice({
  name: "terminal",
  initialState,
  reducers: {
    appendEntry(state, action: PayloadAction<TerminalEntry>) {
      state.entries.push(action.payload);
      if (action.payload.command) {
        state.history.push(action.payload.command);
      }
    },
    clearEntries(state) {
      state.entries = [
        {
          id: "cleared",
          lines: ["CONSOLE CLEARED", "Type HELP to restore the command index."],
          tone: "muted",
        },
      ];
    },
  },
});

export const { appendEntry, clearEntries } = terminalSlice.actions;

export const makeStore = () =>
  configureStore({
    reducer: {
      terminal: terminalSlice.reducer,
    },
  });

export type AppStore = ReturnType<typeof makeStore>;
export type RootState = ReturnType<AppStore["getState"]>;
export type AppDispatch = AppStore["dispatch"];
