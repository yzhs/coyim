package definitions

func init() {
	add(`AddContact2`, &defAddContact2{})
}

type defAddContact2 struct{}

func (*defAddContact2) String() string {
	return `<interface>
  <object class="GtkListStore" id="accounts-model">
    <columns>
      <!-- account name -->
      <column type="gchararray"/>
      <!-- account id -->
      <column type="gchararray"/>
    </columns>
  </object>

  <object class="GtkTextBuffer" id="subscriptionAskMessage">
    <property name="text" translatable="yes">I would like to add you to my contact list.</property>
  </object>

  <object class="GtkAssistant" id="AddContact">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="border_width">6</property>
    <property name="title" translatable="yes">Add new contact</property>
    <property name="resizable">True</property>
    <property name="default-height">400</property>
    <property name="default-width">500</property>
    <property name="destroy-with-parent">true</property>

    <signal name="close" handler="on_cancel_signal" />
    <signal name="cancel" handler="on_cancel_signal" />
    <signal name="escape" handler="on_escape_signal" />
    <signal name="apply" handler="on_apply_signal" />
    <signal name="prepare" handler="on_prepare_signal" />
    <child>
        <object class="GtkBox" id="intro">
            <property name="homogeneous">false</property>
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
            <property name="spacing">6</property>
            <child>
                <object class="GtkBox" id="notification-area">
                    <property name="visible">true</property>
                    <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
                </object>
                <packing>
                    <property name="expand">false</property>
                    <property name="fill">true</property>
                    <property name="position">0</property>
                </packing>
            </child>

            <child>
                <object class="GtkGrid" id="grid">
                    <property name="margin-top">15</property>
                    <property name="margin-bottom">10</property>
                    <property name="margin-start">10</property>
                    <property name="margin-end">10</property>
                    <property name="row-spacing">12</property>
                    <property name="column-spacing">6</property>

                    <child>
                        <object class="GtkLabel" id="subscriptionInstructions">
                            <property name="label" translatable="yes">Please fill in the information for the peer you would like to add. If you have more than one account, ensure you choose the account you would like to add this peer to. You can optionally choose to save a nickname for this peer. You can also choose to automatically allow the peer to see you, without them having to ask permission first. Finally, you can also customize the message the peer will see when asking to see them. Remember that this message will not be end-to-end encrypted, so do not reveal any sensitive information in this message.</property>
                            <property name="visible">true</property>
                            <property name="wrap">true</property>
                            <property name="max-width-chars">65</property>
                        </object>
                        <packing>
                            <property name="left-attach">0</property>
                            <property name="top-attach">0</property>
                            <property name="width">2</property>
                        </packing>
                    </child>

                    <child>
                        <object class="GtkLabel" id="accountsLabel" >
                            <property name="label" translatable="yes">Account:</property>
                            <property name="justify">GTK_JUSTIFY_RIGHT</property>
                            <property name="halign">GTK_ALIGN_END</property>
                        </object>
                        <packing>
                            <property name="left-attach">0</property>
                            <property name="top-attach">1</property>
                        </packing>
                    </child>
                    <child>
                        <object class="GtkComboBox" id="accounts">
                            <property name="model">accounts-model</property>
                            <property name="has-focus">true</property>
                            <property name="hexpand">True</property>
                            <child>
                                <object class="GtkCellRendererText" id="account-name-rendered"/>
                                <attributes>
                                    <attribute name="text">0</attribute>
                                </attributes>
                            </child>
                        </object>
                        <packing>
                            <property name="left-attach">1</property>
                            <property name="top-attach">1</property>
                        </packing>
                    </child>

                    <child>
                        <object class="GtkLabel" id="accountLabel" >
                            <property name="label" translatable="yes">Contact to add:</property>
                            <property name="justify">GTK_JUSTIFY_RIGHT</property>
                            <property name="halign">GTK_ALIGN_END</property>
                        </object>
                        <packing>
                            <property name="left-attach">0</property>
                            <property name="top-attach">2</property>
                        </packing>
                    </child>
                    <child>
                        <object class="GtkEntry" id="address">
                            <property name="placeholder-text">someone@jabber.org</property>
                            <property name="hexpand">True</property>
                            <signal name="activate" handler="on_save_signal" />
                        </object>
                        <packing>
                            <property name="left-attach">1</property>
                            <property name="top-attach">2</property>
                        </packing>
                    </child>

                    <child>
                        <object class="GtkLabel" id="nicknameLabel" >
                            <property name="label" translatable="yes">Nickname:</property>
                            <property name="justify">GTK_JUSTIFY_RIGHT</property>
                            <property name="halign">GTK_ALIGN_END</property>
                        </object>
                        <packing>
                            <property name="left-attach">0</property>
                            <property name="top-attach">3</property>
                        </packing>
                    </child>
                    <child>
                        <object class="GtkEntry" id="nickname">
                            <property name="hexpand">True</property>
                            <signal name="activate" handler="on_save_signal" />
                        </object>
                        <packing>
                            <property name="left-attach">1</property>
                            <property name="top-attach">3</property>
                        </packing>
                    </child>

                </object>
                <packing>
                    <property name="expand">true</property>
                    <property name="fill">true</property>
                    <property name="position">1</property>
                </packing>
            </child>

            <child>
                <object class="GtkCheckButton" id="auto_authorize_checkbutton">
                    <property name="label" translatable="yes">A_llow this contact to view my status</property>
                    <property name="use_action_appearance">False</property>
                    <property name="visible">True</property>
                    <property name="can_focus">True</property>
                    <property name="receives_default">False</property>
                    <property name="no_show_all">True</property>
                    <property name="use_underline">True</property>
                    <property name="xalign">0</property>
                    <property name="active">True</property>
                    <property name="draw_indicator">True</property>
                </object>
                <packing>
                    <property name="expand">False</property>
                    <property name="fill">False</property>
                    <property name="position">2</property>
                </packing>
            </child>
            <child>
                <object class="GtkScrolledWindow" id="message_scrolledwindow">
                    <property name="visible">True</property>
                    <property name="can_focus">True</property>
                    <property name="no_show_all">True</property>
                    <property name="border_width">6</property>
                    <property name="shadow_type">etched-in</property>
                    <property name="min_content_height">5</property>
                    <child>
                        <object class="GtkTextView" id="message_textview">
                            <property name="visible">True</property>
                            <property name="can_focus">True</property>
                            <property name="wrap_mode">word</property>
                            <property name="buffer">subscriptionAskMessage</property>
                        </object>
                    </child>
                </object>
                <packing>
                    <property name="expand">True</property>
                    <property name="fill">True</property>
                    <property name="position">3</property>
                </packing>
            </child>

            <!-- <child internal&#45;child="action_area"> -->
            <!--     <object class="GtkButtonBox" id="button_box"> -->
            <!--         <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property> -->
            <!--         <child> -->
            <!--             <object class="GtkButton" id="button_cancel"> -->
            <!--                 <property name="label">_Cancel</property> -->
            <!--                 <property name="use&#45;underline">True</property> -->
            <!--                 <signal name="clicked" handler="on_close_signal" /> -->
            <!--             </object> -->
            <!--         </child> -->
            <!--         <child> -->
            <!--             <object class="GtkButton" id="button_ok"> -->
            <!--                 <property name="label" translatable="yes">Add</property> -->
            <!--                 <property name="use&#45;underline">True</property> -->
            <!--                 <property name="can&#45;default">true</property> -->
            <!--                 <signal name="clicked" handler="on_next_signal" /> -->
            <!--             </object> -->
            <!--         </child> -->
            <!--     </object> -->
            <!-- </child> -->
        </object>
        <packing>
            <property name="page_type">intro</property>
        </packing>
    </child>
    <child>
        <object class="GtkLabel" id="label2">
            <property name="visible">True</property>
            <property name="label" translatable="yes">Content page</property>
        </object>
        <packing>
            <property name="page_type">content</property>
        </packing>
    </child>
    <child>
        <object class="GtkLabel" id="label3">
            <property name="visible">True</property>
            <property name="label" translatable="yes">Confirmation page</property>
        </object>
        <packing>
            <property name="page_type">confirm</property>
        </packing>
    </child>
  </object>
</interface>
`
}