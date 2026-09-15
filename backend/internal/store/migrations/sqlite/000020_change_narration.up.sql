-- Anomaly narration (issue #238, piece 2/3): cache for the on-demand
-- plain-English "explain this change" AI narration, so repeat views reuse
-- the stored text instead of re-calling the AI provider.
ALTER TABLE changes ADD COLUMN narration TEXT NOT NULL DEFAULT '';
