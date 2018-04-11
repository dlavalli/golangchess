from smshubbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Smshubbeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        smshubbeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("smshubbeat is running"))
        exit_code = smshubbeat_proc.kill_and_wait()
        assert exit_code == 0
