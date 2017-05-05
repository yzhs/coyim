package definitions

func init() {
	add(`StartVerificationNotification`, &defStartVerificationNotification{})
}

type defStartVerificationNotification struct{}

func (*defStartVerificationNotification) String() string {
	return `<interface>
  <object class="GtkInfoBar" id="infobar">
    <property name="message-type">GTK_MESSAGE_WARNING</property>
    <property name="show-close-button">true</property>
    <signal name="response" handler="close_verification"/>
    <child internal-child="content_area">
      <object class="GtkBox" id="box">
        <property name="homogeneous">false</property>
        <child>
        <object class="GtkGrid">
          <property name="column-spacing">4</property>
          <property name="row-spacing">10</property>
          <child>
          <object class="GtkBox">
            <child>
              <object class="GtkImage">
                <property name="file">build/images/alert.png</property>
                <property name="margin-start">8</property>
                <property name="margin-end">5</property>
              </object>
            </child>
            <child>
              <object class="GtkLabel" id="message">
                <property name="wrap">true</property>
                <property name="justify">GTK_JUSTIFY_LEFT</property>
              </object>
            </child>
          </object>
          <packing>
            <property name="left-attach">0</property>
            <property name="top-attach">0</property>
          </packing>
          </child>
          <child>
            <object class="GtkButton" id="button_verify">
              <property name="label" translatable="yes">Validate Secure Channel</property>
              <property name="can-default">true</property>
            </object>
            <packing>
              <property name="left-attach">0</property>
              <property name="top-attach">1</property>
            </packing>
          </child>
        </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}
