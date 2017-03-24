package definitions

func init() {
	add(`OopsWrongPIN`, &defOopsWrongPIN{})
}

type defOopsWrongPIN struct{}

func (*defOopsWrongPIN) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="box">
        <property name="border-width">10</property>
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkGrid">
            <property name="halign">GTK_ALIGN_CENTER</property>
            <property name="row-spacing">60</property>
            <child>
              <object  class="GtkImage">
                <property name="file">build/images/red_alert.png</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel">
                <property name="label" translatable="yes">WRONG PIN!</property>
              </object>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">0</property>
              </packing>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkLabel">
            <property name="label" translatable="yes">You may have typed it wrong or this channel is insecure. You can validate it again later.</property>
          </object>
        </child>
        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="continue_anyway">
                <property name="label" translatable="yes">Continue anyway</property>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="get_out">
                <property name="label" translatable="yes">Get me out of here</property>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
    <action-widgets>
      <action-widget response="no" default="true">continue_anyway</action-widget>
      <action-widget response="yes">get_out</action-widget>
    </action-widgets>
  </object>
</interface>
`
}
