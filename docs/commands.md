# Commands

Commands are executed at 3 specific points:

#### pre_commands
Pre-commands are used to prepare the environment for execution, OS detection and setting of runtime flags is done in this phase so that they can be used in all other phases. e.g. set an environment variable based on the output of a command.

#### commands
Phases can only append to this Commands list.

#### post_commands
Post commands run after all the phases have completed and can be used for cleanup functions are for handing off to other systems.
